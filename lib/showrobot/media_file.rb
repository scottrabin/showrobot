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

		def is_movie?
			parse[:type] == :movie
		end

		def is_tv?
			parse[:type] == :tv
		end

		def season
			parse[:season]
		end

		def episode
			parse[:episode]
		end

		def series database
			ShowRobot.datasource_for(database).new(self).series
		end

		def name_guess
			@name_guess ||= parse[:name_guess].gsub(/[^a-zA-Z0-9]/, ' ').gsub(/\s+/, ' ').strip
		end

		protected
		def initialize fileName
			@fileName = fileName
		end

		# parses a file name for the constituent parts
		def parse
			@parse ||= ShowRobot.parse_filename File.basename(@fileName)
		end

		private
		@@video_types = {}
	end

end
