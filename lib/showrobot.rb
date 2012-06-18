module ShowRobot

end

require 'showrobot/media_file'
(
	# load all utility files
	Dir[File.dirname(__FILE__) + '/showrobot/utility/*.rb'] +
	# load all media file parsers
	Dir[File.dirname(__FILE__) + '/showrobot/video/*.rb']
).each { |file| require file }
