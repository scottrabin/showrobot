require 'text'

def by_distance to, prop = nil
	if prop.nil?
		Proc.new do |a, b|
			Text::Levenshtein.distance(a, to) - Text::Levenshtein.distance(b, to)
		end
	else
		Proc.new do |a, b|
			Text::Levenshtein.distance(a[prop], to) - Text::Levenshtein.distance(b[prop], to)
		end
	end
end
