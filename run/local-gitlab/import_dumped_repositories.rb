require 'securerandom'
require 'fileutils'

# === CONFIGURATION ===
base_path = "/repos-on-host"  # Make sure this directory exists and is mounted inside the container
gitlab_repo_base_path = "/var/opt/gitlab/git-data/repositories"

# === VALIDATION ===

raise "Directory '#{base_path}' not found!" unless Dir.exist?(base_path)

organization = ::Organizations::Organization.find_by(name: 'Default')
raise "Organization 'Default' not found!" unless organization

root_user = User.find_by_username('root')
raise "Root user 'root' not found!" unless root_user

puts "Using organization '#{organization.name}' and user '#{root_user.name}'"

# === STEP 1: DELETE ALL EXISTING GROUPS AND PROJECTS ===

puts "Deleting all existing groups and their projects..."

Group.all.each do |group|
  begin
    group.projects.each do |project|
      puts "  Deleting project '#{project.full_path}'..."
      project.destroy!
    end
    puts "Deleting group '#{group.full_path}'..."
    group.destroy!
  rescue => e
    puts "Error deleting group '#{group.full_path}': #{e.message}"
  end
end

puts "All existing groups and projects deleted."

# === STEP 2: CREATE GROUPS & IMPORT REPOSITORIES ===

Dir.entries(base_path).select { |entry|
  File.directory?(File.join(base_path, entry)) && !(entry == '.' || entry == '..')
}.each do |group_name|
  group_path = File.join(base_path, group_name)
  puts "\nProcessing group: #{group_name}"

  # Create group if it doesn't exist
  group = Group.find_by(path: group_name)
  unless group
    group = Group.new(
      name: group_name,
      path: group_name,
      visibility_level: Gitlab::VisibilityLevel::PUBLIC,
      organization: organization
    )
    if group.save
      group.add_owner(root_user)
      puts "  Group '#{group_name}' created."
    else
      puts "  Failed to create group '#{group_name}': #{group.errors.full_messages.join(', ')}"
      next
    end
  end

  # Process each .git repo in group directory
  Dir.entries(group_path).select { |entry| entry.end_with?('.git') }.each do |repo_dir|
    repo_name = File.basename(repo_dir, '.git')
    puts "  Processing repository: #{repo_name}"

    project = Project.find_by(path: repo_name, namespace: group)
    if project
      puts "    Project '#{repo_name}' already exists. Skipping."
      next
    end

    begin
      project = Project.new(
        name: repo_name,
        path: repo_name,
        namespace: group,
        creator: root_user,
        visibility_level: Gitlab::VisibilityLevel::PRIVATE,
        organization_id: organization.id
      )

      project.save!

      # Prepare paths
      repo_source_path = File.join(group_path, repo_dir)
      repo_destination_path = File.join(gitlab_repo_base_path, project.disk_path + '.git')

      # Remove the auto-created empty repo
      if Dir.exist?(repo_destination_path)
        puts "    Removing existing repository at '#{repo_destination_path}'..."
        FileUtils.rm_rf(repo_destination_path)
      end

      # Copy your actual bare repo
      FileUtils.mkdir_p(File.dirname(repo_destination_path))
      FileUtils.cp_r(repo_source_path, repo_destination_path)
      FileUtils.chown_R('git', 'git', repo_destination_path)

      puts "    Repository '#{repo_name}' successfully imported."
    rescue => e
      puts "    Failed to import repository '#{repo_name}': #{e.message}"
      next
    end
  end
end

puts "\n✅ All groups and repositories processed successfully."
