module ShowRobot

	class MediaFile
		# class methods
		def self.load fileName
			begin
				@@video_types[File.extname(fileName)].new fileName
			rescue
				raise "No parser exists for files of type '#{File.extname(fileName)}'"
			end
		end

		def self.isvideo? fileName
			@@video_types.include? File.extname(fileName)
		end

		def self.addType ext, klass
			@@video_types['.' + ext.to_s] = klass
		end

		# instance methods
		def isvideo?
			MediaFile.isvideo? @fileName
		end

		def season
			self.parse
			@season
		end

		def episode
			self.parse
			@episode
		end

		attr_accessor :nameGuess

		protected
		def initialize fileName
			@fileName = File.basename(fileName)
			@nameGuess = @fileName.dup
		end

		# parses a file name for the constituent parts
		def parse
			if @parse.nil?
				@season, @episode, matchStart, matchLength = ShowRobot.get_season_episode(@fileName)
				@nameGuess[matchStart, matchLength] = ''
				@parse = true
			end
			[@season, @episode]
		end

		private
		@@video_types = {}
	end

end
