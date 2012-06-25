module ShowRobot

	class Datasource
	end

	class << self
		DATASOURCES = {}
		def add_datasource sym, klass
			DATASOURCES[sym] = klass
		end

		def datasource_for sym
			DATASOURCES[sym.to_sym]
		end

		def url_encode(s)
			s.to_s.gsub(/[^a-zA-Z0-9_\-.]/n){ sprintf("%%%02X", $&.unpack("C")[0]) }
		end
	end

end
