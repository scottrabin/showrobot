require 'helper'

describe ShowRobot, "name identification" do
	with_file 'test 1x02.avi' do
		its(:name_guess) { should eq('test') }
	end
end

describe ShowRobot, "season/episode identification" do
	with_file 'test.S01E02.avi' do
		its(:season) { should eq(1) }
		its(:episode) { should eq(2) }
	end

	with_file 'example.01x02.avi' do
		its(:season) { should eq(1) }
		its(:episode) { should eq(2) }
	end
end
