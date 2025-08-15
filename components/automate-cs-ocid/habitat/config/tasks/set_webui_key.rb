require 'json'
require 'fileutils'

webui_key = ENV['WEBUI_KEY']
private_chef_secrets_path = ENV['PRIVATE_CHEF_SECRETS_PATH']

# Create the private-chef-secrets.json file if it doesn't exist
unless File.exist?(private_chef_secrets_path)
  puts "Creating private-chef-secrets.json file at #{private_chef_secrets_path}"
  # Create directory if it doesn't exist
  FileUtils.mkdir_p(File.dirname(private_chef_secrets_path))
  # Create initial empty JSON structure
  initial_secrets = {
    "chef-server" => {}
  }
  File.write(private_chef_secrets_path, JSON.pretty_generate(initial_secrets))
  puts "Successfully created private-chef-secrets.json"
end

# Verify the file exists before trying to read it
unless File.exist?(private_chef_secrets_path)
  puts "ERROR: Failed to create or find private-chef-secrets.json at #{private_chef_secrets_path}"
  exit 1
end

private_chef_secrets_content = File.read(private_chef_secrets_path)
private_chef_secrets_hash = JSON.parse(private_chef_secrets_content)

# Ensure chef-server key exists
private_chef_secrets_hash["chef-server"] ||= {}
private_chef_secrets_hash["chef-server"]["webui_key"] = webui_key

File.write(private_chef_secrets_path, JSON.pretty_generate(private_chef_secrets_hash))
