# test the bin/showrobot command line executable
require 'helper'

describe ShowRobot, 'command line executable' do

	# TODO - configuration editing tests
	
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
				%x(#{File.dirname(__FILE__)}/../bin/showrobot rename -Dm "/tmp/showrobot_test_2.avi" /tmp/showrobot_test.avi)

				File.exists?('/tmp/showrobot_test.avi').should be(true)
				File.exists?('/tmp/showrobot_test_2.avi').should be(false)
				File.delete('/tmp/showrobot_test.avi')
			end
		end

		describe '--movie-database' do
			it 'should use the specified database for movies' do
				output = `./bin/showrobot rename -Dv --movie-database mockmovie -m "nil" "#{MOVIE_NAME}"`
				# make sure the output contains a line like "from Some DB"
				output.should include('from Mock Movie Database')
			end
		end

		describe '--tv-database' do
		end

		describe '--prompt' do
		end

		describe '--force-movie' do
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
