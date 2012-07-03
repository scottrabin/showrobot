module ShowRobot

	def self.log level, message, opts = {}
		puts "[#{level}] #{message}" if ShowRobot.config[:verbose]
		if not opts[:no_write]
			File.open(config[:basepath] + '/log', 'a') do |file|
				file.printf "[%s](%-8s) %s\n", Time.new.strftime("%y/%m/%d %H:%M:%S"), level, message
			end
		end
	end

end
