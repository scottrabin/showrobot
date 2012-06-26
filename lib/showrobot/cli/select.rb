# helper function for prompts
def select list, prompt, &format
	range = 0...10
	begin
		puts prompt
		range.each do |i|
			puts format.call i, list[i]
		end
		print " > "
		input = $stdin.gets.chomp

		case input
		when /^[\d]+$/
			# it's a number, make sure it's in the range
			selected = list[input.to_i] if input.to_i < list.length
		when 'n' # next range
			range = Range.new(range.end, [range.end + 10, list.length].min, true) if range.end < list.length
		when 'p' # prev range
			range = Range.new([0, range.first - 10].max, range.first, true) if range.first > 1
		when 'q'
			exit
		end
	end while selected.nil?
	selected
end
