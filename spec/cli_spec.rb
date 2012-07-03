# test the bin/showrobot command line executable
require 'helper'

describe ShowRobot, 'command line executable' do

	# TODO - configuration editing tests
	
	CLI = File.expand_path(File.dirname(__FILE__) + '/../bin/showrobot')
	MOVIE_NAME = 'First Movie (2000).avi'
	TV_SHOW_NAME = 'ShowSeries.S01E02.TheTitle.avi'

	describe 'when running: showrobot rename' do
		# TODO --config-file, -f
		# TODO --verbose, -v
		# TODO --cache-dir, -c
		# TODO --no-cache, -n
		# TODO --clear-cache, -C

		describe '--dry-run' do
			it 'should not actually move the file' do
				`touch /tmp/showrobot_test.avi`
				`#{CLI} rename -Dm "/tmp/showrobot_test_2.avi" /tmp/showrobot_test.avi`

				File.exists?('/tmp/showrobot_test.avi').should be(true)
				File.exists?('/tmp/showrobot_test_2.avi').should be(false)
				File.delete('/tmp/showrobot_test.avi')
				$?.should eq(0)
			end
		end

		describe '--movie-database' do
			it 'should use the specified database for movies' do
				output = cli %(rename -Dv --movie-database mockmovie "#{MOVIE_NAME}")
				# make sure the output contains a line like "from Some DB"
				output.should include("from #{ShowRobot::MockMovie::DB_NAME}")
			end
		end

		describe '--tv-database' do
			it 'should use the specified database for tv shows' do
				output = cli %(rename -Dv --tv-database mocktv "#{TV_SHOW_NAME}")
				output.should include("from #{ShowRobot::MockTV::DB_NAME}")
			end
		end

		describe '--prompt' do
			# TODO - not sure how to test interactive prompts
		end

		describe '--force-movie' do
			it 'should use the movie database even though it looks like a TV show' do
				# make sure the TV episode actually gets parsed as a TV episode automatically
				verify = cli %(rename -Dv --tv-database mocktv "#{TV_SHOW_NAME}")
				verify.should include("from #{ShowRobot::MockTV::DB_NAME}")

				# now make sure it gets forced to the mock movie db with --force-movie
				movie = cli %(rename -Dv --movie-database mockmovie --force-movie "#{TV_SHOW_NAME}")
				movie.should include("from #{ShowRobot::MockMovie::DB_NAME}")
			end
		end

		describe '--force-tv' do
			it 'should use the tv database even though it looks like a movie' do
				# make sure the movie actually gets parsed as a movie automatically
				verify = cli %(rename -Dv --movie-database mockmovie "#{MOVIE_NAME}")
				verify.should include("from #{ShowRobot::MockMovie::DB_NAME}")

				# now make sure it gets forced to the mock movie db with --force-movie
				movie = cli %(rename -Dv --tv-database mocktv --force-tv "#{MOVIE_NAME}")
				movie.should include("from #{ShowRobot::MockTV::DB_NAME}")
			end
		end

		describe '--movie-format' do
			it 'should format the output parameters correctly' do
				# for movies, {t} => title, {y} => year, {ext} => extension
				title = cli %(rename -Dv --movie-database mockmovie -m "{t}" "#{MOVIE_NAME}")
				File.basename(title[ShowRobotHelper::CMD_RENAME_TO, 1]).should eq('First Movie')

				title = cli %(rename -Dv --movie-database mockmovie -m "{y}" "#{MOVIE_NAME}")
				File.basename(title[ShowRobotHelper::CMD_RENAME_TO, 1]).should eq('2000')

				title = cli %(rename -Dv --movie-database mockmovie -m "{ext}" "#{MOVIE_NAME}")
				File.basename(title[ShowRobotHelper::CMD_RENAME_TO, 1]).should eq('avi')

			end
		end

		describe '--tv-format' do
		end

		describe '--space-replace' do
		end

	end
end
