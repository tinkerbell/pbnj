require 'grpc'
require 'api/v1/machine_services_pb'
require "pbnj/task"


module PBnJ
    class Client
        attr_accessor :pbnj_ip, :pbnj_port

        def initialize(config = {})
            config.each_pair do |option, value|
                send("#{option}=", value)
            end
            @v1 = Pbnj::Api::V1
        end

        def stub_connect(conn)
            return conn.new("#{pbnj_ip}:#{pbnj_port}", :this_channel_is_insecure)
        end

        def power_status(ip, user, pass)
            stub = stub_connect(@v1::Machine::Stub)
            power_request = @v1::PowerRequest.new(
                authn: @v1::Authn.new(
                    directAuthn: @v1::DirectAuthn.new(
                        host: @v1::Host.new(
                            host: ip
                        ),
                        username: user,
                        password: pass
                    )
                ),
                power_action: @v1::PowerAction::POWER_ACTION_STATUS
            )
            task_from_response(stub.power(power_request).task_id)
        end

        private

        def task_from_response(task_id)
            t = Task.new(task_id, self)
            return t.wait
        end
    end
end
