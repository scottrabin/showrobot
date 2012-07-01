require 'text'

def distance_differential a, b, to
	Text::Levenshtein.distance(a, to) - Text::Levenshtein.distance(b, to)
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
