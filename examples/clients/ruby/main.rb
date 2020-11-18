this_dir = File.expand_path(File.dirname(__FILE__))
lib_dir = File.join(this_dir, 'lib')
$LOAD_PATH.unshift(lib_dir) unless $LOAD_PATH.include?(lib_dir)

require 'pbnj/client'
require 'json'

include PBnJ

def main
    host = ARGV.shift || "localhost"
    user = ARGV.shift || "ADMIN"
    pass = ARGV.shift || "ADMIN"

    config = {"pbnj_ip"=> "localhost", "pbnj_port"=> "9090"}
    client = PBnJ::Client.new(config)
    task_response = client.power_status(host, user, pass)
    
    jj JSON[task_response.to_json]

end

main
