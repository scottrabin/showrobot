module ShowRobot

	class AVIFile < MediaFile
		def initialize fileName
			super(fileName)
		end

		def duration
			@duration ||= `ffmpeg -i "#{@fileName}" 2>&1`[/Duration: ([\d:\.]*)/, 1].split(':').each_with_index.map { |n, i| n.to_f * (60 ** (2-i)) }.reduce(0, :+) rescue nil
		end

	end

	MediaFile.addType :avi, AVIFile

end
