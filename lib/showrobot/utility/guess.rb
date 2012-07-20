module ShowRobot

	FEATURE_MU = [0, 617.257452, 1.072241, 0.388629, 76.056187]
	FEATURE_SIGMA = [1, 439.334486, 0.984552, 0.487602, 43.708610]
	FEATURE_THETA = [3.824134, 0.917716, 1.274485, -0.204749, -13.159588]

	def compute_z feat
		# intercept
		feat.unshift 1
		(0...feat.length).reduce(0) do |memo, index|
			memo + FEATURE_THETA[index] * ((feat[index] - FEATURE_MU[index])/FEATURE_SIGMA[index])
		end
	end

	def features mediaFile
		file = mediaFile.fileName
		[
			# feature 1: size of file, in MB
			File.size(file) / 1048576.0,
			# feature 2: no. 2 digit numbers
			file.scan(/(?<=\D)\d{2}(?=\D)/).length,
			# feature 3: presence of 4-digit number
			file[/\d{4}/] ? 1 : 0,
			# feature 4: duration, in minutes
			mediaFile.duration / 60.0
		]
	end

	def sigmoid val
		1/(1 + Math.exp(-val))
	end

	def guess file
		if sigmoid(compute_z(features file)) < 0.5
			:movie
		else
			:tv
		end
	end

	module_function :compute_z, :features, :sigmoid, :guess

end
