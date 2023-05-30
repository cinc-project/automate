require 'securerandom'

begin
  secret_key_base = SecureRandom.hex(64)
  
  ocid_config_folder_path = ENV['OCID_CONFIG_FOLDER_PATH']
  secret_key_base_file_path = ENV['SECRET_KEY_BASE_FILE_PATH']
  
  if File.exists?(secret_key_base_file_path)
    puts "SECRET_KEY_BASE file exists. File will be overwritten.."
  else
    puts "SECRET_KEY_BASE file doesn't exist. It will be generated.."
    dir = File.dirname(secret_key_base_file_path)
    unless File.directory?(dir)
      FileUtils.mkdir_p(dir)
    end
  end
  
  File.write(secret_key_base_file_path, secret_key_base)  
rescue StandardError => e
  puts "ERROR: Failed to generate secret_key_base for OCID."
  puts "SYS ERROR: #{e.inspect}"
end