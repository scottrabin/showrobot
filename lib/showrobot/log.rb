module ShowRobot

	def self.log level, message
		File.open(config[:basepath] + '/log', 'a') do |file|
			file.printf "[%s](%-8s) %s\n", Time.new.strftime("%y/%m/%d %H:%M:%S"), level, message
		end
	end

end
