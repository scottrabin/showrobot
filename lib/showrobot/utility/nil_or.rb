# TODO - make sure there isn't a way in core Ruby to do this same thing as a one-liner
def nil_or method, item
	if not item.nil? and item.respond_to? method.to_sym
		item.send method.to_sym
	end
end
