# test the bin/showrobot command line executable
require 'helper'
require 'fileutils'

describe ShowRobot, 'command line executable' do

	# TODO - configuration editing tests
	
	MOVIE = {
		:filename   => 'Fourth Movie: The Subtitle.(2003).avi',
		:title      => 'Fourth Movie: The Subtitle',
		:runtime    => 92,
		:year       => '2003',
		:extension  => 'avi'
	}
	TV_SHOW = {
		:filename   => 'AnotherSeries.S01E02.TheTitle.avi',
		:seriesname => 'Another Series',
		:season     => '01',
		:episode    => '02',
		:title      => 'Fifth Episode',
		:extension  => 'avi'
	}


	describe 'when running: showrobot rename' do
		# TODO --config-file, -f
		# TODO --verbose, -v
		# TODO --cache-dir, -c
		# TODO --no-cache, -n
		# TODO --clear-cache, -C

		describe '--dry-run' do
			before :each do
				@startfile = temp_file 'showrobot_test.avi'
				@outfile   = temp_file 'failure.avi'
				File.delete @startfile if File.exists?(@startfile)
				File.delete @outfile   if File.exists?(@outfile)
				FileUtils.mkdir_p File.dirname(@startfile)
				FileUtils.touch @startfile
			end
			
			after :each do
				File.delete @startfile if File.exists?(@startfile)
				File.delete @outfile   if File.exists?(@outfile)
			end

			it 'should not actually move the file' do
				move = cli %(rename -Dm "#{@outfile}" "#{@startfile}")
				File.exists?(@startfile).should == true
				File.exists?(@outfile).should == false
			end

			it 'should move the file without the option' do
				move = cli %(rename -m "#{@outfile}" "#{@startfile}")
				File.exists?(@startfile).should == false
				File.exists?(@outfile).should == true
			end
		end

		describe '--movie-database' do
			it 'should use the specified database for movies' do
				output = cli %(rename -Dv --movie-database mockmovie "#{MOVIE[:filename]}")
				# make sure the output contains a line like "from Some DB"
				output.should include("from #{ShowRobot::MockMovie::DB_NAME}")
			end
		end

		describe '--tv-database' do
			it 'should use the specified database for tv shows' do
				output = cli %(rename -Dv --tv-database mocktv "#{TV_SHOW[:filename]}")
				output.should include("from #{ShowRobot::MockTV::DB_NAME}")
			end
		end

		describe '--prompt' do
			# TODO - not sure how to test interactive prompts
		end

		describe '--force-movie' do
			it 'should use the movie database even though it looks like a TV show' do
				# make sure the TV episode actually gets parsed as a TV episode automatically
				verify = cli %(rename -Dv --tv-database mocktv "#{TV_SHOW[:filename]}")
				verify.should include("from #{ShowRobot::MockTV::DB_NAME}")

				# now make sure it gets forced to the mock movie db with --force-movie
				movie = cli %(rename -Dv --movie-database mockmovie --force-movie "#{TV_SHOW[:filename]}")
				movie.should include("from #{ShowRobot::MockMovie::DB_NAME}")
			end
		end

		describe '--force-tv' do
			it 'should use the tv database even though it looks like a movie' do
				# make sure the movie actually gets parsed as a movie automatically
				verify = cli %(rename -Dv --movie-database mockmovie "#{MOVIE[:filename]}")
				verify.should include("from #{ShowRobot::MockMovie::DB_NAME}")

				# now make sure it gets forced to the mock movie db with --force-movie
				movie = cli %(rename -Dv --tv-database mocktv --force-tv "#{MOVIE[:filename]}")
				movie.should include("from #{ShowRobot::MockTV::DB_NAME}")
			end
		end

		describe '--movie-format' do
			it 'should format the output parameters correctly' do
				# for movies, {t} => title, {y} => year, {ext} => extension
				title = cli %(rename -Dv --movie-database mockmovie -m "{t}" "#{MOVIE[:filename]}")
				File.basename(out_file title).should == MOVIE[:title]

				year = cli %(rename -Dv --movie-database mockmovie -m "{y}" "#{MOVIE[:filename]}")
				File.basename(out_file year).should == MOVIE[:year]

				extension = cli %(rename -Dv --movie-database mockmovie -m "{ext}" "#{MOVIE[:filename]}")
				File.basename(out_file extension).should == MOVIE[:extension]
			end
		end

		describe '--tv-format' do
			it 'should format the output parameters correctly' do
				# for tv shows, {n} => series name, {s} => season, {e} => episode,
				#               {t} => episode title, {ext} => file extension
				seriesname = cli %(rename -Dv --tv-database mocktv -t "{n}" "#{TV_SHOW[:filename]}")
				File.basename(out_file seriesname).should == TV_SHOW[:seriesname]

				title = cli %(rename -Dv --tv-database mocktv -t "{t}" "#{TV_SHOW[:filename]}")
				File.basename(out_file title).should == TV_SHOW[:title]

				season = cli %(rename -Dv --tv-database mocktv -t "{s}" "#{TV_SHOW[:filename]}")
				File.basename(out_file season).should == TV_SHOW[:season]

				episode = cli %(rename -Dv --tv-database mocktv -t "{e}" "#{TV_SHOW[:filename]}")
				File.basename(out_file episode).should == TV_SHOW[:episode]

				extension = cli %(rename -Dv --tv-database mocktv -t "{ext}" "#{TV_SHOW[:filename]}")
				File.basename(out_file extension).should == TV_SHOW[:extension]
			end
		end

		describe '--space-replace' do
			it 'should format the output with all spaces in the filename replaced' do
				output = cli %(rename -Dv --tv-database mocktv -t "{n}!{t}" --space-replace "?" "#{TV_SHOW[:filename]}")
				File.basename(out_file output).should == "#{TV_SHOW[:seriesname]}!#{TV_SHOW[:title]}".gsub(/\s/, '?')
			end

			it 'should format the output with ONLY the filename spaces replaced' do
				dir = temp_file 'a b c dir'
				output = cli %(rename -Dv --tv-database mocktv -t "#{dir}/{t}.{ext}" --space-replace "?" "#{TV_SHOW[:filename]}")
				out_file(output).should eq("#{dir}/#{TV_SHOW[:title].gsub(/\s/, '?')}.#{TV_SHOW[:extension]}")
			end
		end

	end
end
