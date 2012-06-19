module ShowRobot

	YEAR_PATTERNS = [
		# in parens
		/\((\d{4})\)/,
		# in square brackets
		/\[(\d{4})\]/,
		# ... really, just any 4 numbers in a row.
		/\D(\d{4})\D/
	]

	def file_get_year fileName
		if not YEAR_PATTERNS.find { |pattern| pattern.match fileName }.nil?
			$1
		end
	end

	module_function :file_get_year

end
