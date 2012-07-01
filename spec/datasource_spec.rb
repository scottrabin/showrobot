require 'helper'
require 'mock/tvsource'

describe ShowRobot, 'datasource API' do

	mockdb = ShowRobot.create_datasource 'mocktv'
	mockdb.mediaFile = ShowRobot::MediaFile.load 'SeriesName.S02E01.avi'

	describe 'when querying for the season list' do
		it 'should return a valid list of series' do
			mockdb.series_list.should have_at_least(15).items
			mockdb.series_list.each do |series|
				series[:name].should be_an_instance_of(String)
				series[:source].should be_an_instance_of(Hash)
			end
		end
	end
end
