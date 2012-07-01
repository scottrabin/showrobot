def translate str, mapping
	str.gsub(/(?<!\\){[^}]+}/, mapping).gsub(/\\({[^}]+})/, '\1')
end
