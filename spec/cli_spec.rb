# test the bin/showrobot command line executable
require 'helper'

describe ShowRobot, 'command line executable' do

	# TODO - configuration editing tests
	
	CLI = File.expand_path(File.dirname(__FILE__) + '/../bin/showrobot')
	MOVIE_NAME = 'A_Movie (2013).avi'
	TV_SHOW_NAME = 'ShowSeries.S01E02.TheTitle.avi'

	describe 'when running: showrobot rename' do
		# TODO --config-file, -f
		# TODO --verbose, -v
		# TODO --cache-dir, -c
		# TODO --no-cache, -n
		# TODO --clear-cache, -C

		describe '--dry-run' do
			it 'should not actually move the file' do
				%x(touch /tmp/showrobot_test.avi)
				%x(#{CLI} rename -Dm "/tmp/showrobot_test_2.avi" /tmp/showrobot_test.avi)

				File.exists?('/tmp/showrobot_test.avi').should be(true)
				File.exists?('/tmp/showrobot_test_2.avi').should be(false)
				File.delete('/tmp/showrobot_test.avi')
			end
		end

		describe '--movie-database' do
			it 'should use the specified database for movies' do
				output = `#{CLI} rename -Dv --movie-database mockmovie -m "nil" "#{MOVIE_NAME}"`
				# make sure the output contains a line like "from Some DB"
				output.should include('from Mock Movie Database')
			end
		end

		describe '--tv-database' do
			it 'should use the specified database for tv shows' do
				output = `#{CLI} rename -Dv --tv-database mocktv -t "nil" "#{TV_SHOW_NAME}"`

				output.should include('from Mock TV Database')
			end
		end

		describe '--prompt' do
			# TODO - not sure how to test interactive prompts
		end

		describe '--force-movie' do
			it 'should use the movie database even though it looks like a TV show' do
				# make sure the TV episode actually gets parsed as a TV episode automatically
				verify = `#{CLI} rename -Dv --tv-database mocktv -t "nil" "#{TV_SHOW_NAME}"`
				verify.should include("from #{ShowRobot::MockTV::DB_NAME}")

				# now make sure it gets forced to the mock movie db with --force-movie
				movie = `#{CLI} rename -Dv --movie-database mockmovie -m "nil" --force-movie "#{TV_SHOW_NAME}"`
				movie.should include("from #{ShowRobot::MockMovie::DB_NAME}")
			end
		end

		describe '--force-tv' do
		end

		describe '--movie-format' do
		end

		describe '--tv-format' do
		end

		describe '--space-replace' do
		end

	end
end
