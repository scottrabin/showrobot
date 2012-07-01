require 'helper'
require 'mock/tvsource'

describe ShowRobot, 'datasource API' do

	mockdb = ShowRobot.create_datasource 'mocktv'
	mockdb.mediaFile = ShowRobot::MediaFile.load 'SeriesName.S02E01.avi'

	describe 'when querying for the season list' do
		it 'should return a valid list of series' do
			mockdb.series_list.should have(2).items
			mockdb.series_list.each do |series|
				series[:name].should be_an_instance_of(String)
				series[:source].should be_an_instance_of(Hash)
			end
		end
	end

	describe 'when querying for the episode list' do
		it 'should return a list of episodes' do
			mockdb.episode_list.should have(15).items
			mockdb.episode_list.each do |episode|
				episode[:series].should be_an_instance_of(String)
				episode[:title].should be_an_instance_of(String)
				episode[:season].should be_a_kind_of(Integer)
				episode[:episode].should be_a_kind_of(Integer)
				episode[:episode_ct].should be_a_kind_of(Integer)
			end
		end
	end
end
