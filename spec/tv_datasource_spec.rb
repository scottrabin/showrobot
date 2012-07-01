require 'helper'
require 'mock/tvsource'

describe ShowRobot, 'datasource API' do

	before :each do
		@file = ShowRobot::MediaFile.load 'SeriesName.S02E01.avi'
		@db   = ShowRobot.create_datasource :mocktv
		@db.mediaFile = @file
	end

	describe 'when querying for the season list' do
		it 'should return a valid list of series' do
			@db.series_list.should have(2).items
			@db.series_list.each do |series|
				series[:name].should be_an_instance_of(String)
				series[:source].should be_an_instance_of(Hash)
			end
		end
	end

	describe 'when querying for the episode list' do
		it 'should return a list of episodes' do
			@db.episode_list.should have(15).items
			@db.episode_list.each do |episode|
				episode[:series].should be_an_instance_of(String)
				episode[:title].should be_an_instance_of(String)
				episode[:season].should be_a_kind_of(Integer)
				episode[:episode].should be_a_kind_of(Integer)
				episode[:combined_ep].should be_a_kind_of(Integer)
			end
		end
	end

	describe 'when querying for a specific episode' do
		it 'should return the corresponding episode' do
			@db.episode(1, 1)[:title].should eq('The Fourth Episode')
			@db.episode(2, 3)[:title].should eq('The Ninth Episode')
			@db.episode(4, 2)[:title].should eq('The Fourteenth Episode')
		end

		it 'should automatically return the episode for the associated media file' do
			@db.episode[:title].should eq('The Seventh Episode')
		end
	end

	describe 'when the series is explicitly specified' do
		it 'should match the correct series' do
			@db.series = 'Another Series'
			@db.episode_list.each do |episode|
				episode[:series].should eq('Another Series')
			end
		end
	end

end
