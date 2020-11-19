require "pbnj/errors"
require 'api/v1/task_services_pb'

module PBnJ
    class Task
      LONG_POLLING_TIMEOUT = 30

      attr_accessor :id, :client
  
      def initialize(id, client)
        self.id = id
        self.client = client
        @v1 = Pbnj::Api::V1
      end
  
      def done?
        entity.complete
      end
  
      def failed?
        done? && errors.any?
      end
  
      def errors
        entity.error.code != 0 || []
      end
  
      def wait
        sleep 1 until done?  
        entity
      end
  
      private
  
      def entity
        $i = 0
        $num = LONG_POLLING_TIMEOUT*2
        stub_status = client.stub_connect(@v1::Task::Stub)
        while $i < $num  do
          begin
          task_request = @v1::StatusRequest.new(task_id: "#{id}")
          task_response = stub_status.status(task_request)
          if task_response.state == 'complete'
            @finished_task = task_response
            break
          end
          rescue GRPC::BadStatus => e
            abort "ERROR: #{e.message}"
          end
          sleep(2)
          $i +=1
        end
        return @finished_task if @finished_task
      end
    end
  end
