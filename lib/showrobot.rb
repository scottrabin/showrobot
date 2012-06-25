module ShowRobot

	class << self
		attr_accessor :verbose
	end

end

require 'showrobot/media_file'
require 'showrobot/db'
(
	# load all utility files
	Dir[File.dirname(__FILE__) + '/showrobot/utility/*.rb'] +
	# load all media file parsers
	Dir[File.dirname(__FILE__) + '/showrobot/video/*.rb'] +
	# load all database shims
	Dir[File.dirname(__FILE__) + '/showrobot/db/*.rb']
).each { |file| require file }
