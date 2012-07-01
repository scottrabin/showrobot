require 'text'

def word_distance a, b
	Text::Levenshtein.distance(a, b)
end

def distance_differential a, b, to
	word_distance(a, to) - word_distance(b, to)
end

def by_distance to, prop = nil
	if prop.nil?
		Proc.new do |a, b|
			distance_differential a, b, to
		end
	else
		Proc.new do |a, b|
			distance_differential a[prop], b[prop], to
		end
	end
end
