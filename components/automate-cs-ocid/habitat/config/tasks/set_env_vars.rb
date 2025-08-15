require 'yaml'
require 'fileutils'

env_config_file_path = ENV['ENV_CONFIG_FILE_PATH']

# Ensure the directory exists
config_dir = File.dirname(env_config_file_path)
FileUtils.mkdir_p(config_dir) unless Dir.exist?(config_dir)

# Check if the file exists before trying to load it
if File.exist?(env_config_file_path)
  env_config = YAML.load_file(env_config_file_path)
else
  puts "Warning: Environment config file not found at #{env_config_file_path}. Creating new config."
  env_config = {}
end
override_configs = {
  'endpoint' => ENV['CHEF_SERVER_ENDPOINT'],
  'superuser' => ENV['CHEF_SERVER_SUPERUSER'],
  'ssl_verify_mode' => ENV['CHEF_SERVER_SSL_VERIFY_MODE']
}
env_config['chef'] ||= {}
env_config['chef'].merge!(override_configs)

# Overriding the default `secret_key_base` with the custom/generated value
# of `secret_key_base` set in the environment only if the env variable is set.
env_config['secret_key_base'] = ENV['SECRET_KEY_BASE'] unless (ENV['SECRET_KEY_BASE'].nil? || ENV['SECRET_KEY_BASE'].empty?)

env_config_yml = env_config.to_yaml
File.write(env_config_file_path, env_config_yml)
