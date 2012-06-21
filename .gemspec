# -*- encoding: utf-8 -*-
lib = File.expand_path('../lib/', __FILE__)
$:.unshift lib unless $:.include?(lib)

require 'showrobot/version'

Gem::Specification.new do |s|
	s.name          = 'showrobot'
	s.version       = ShowRobot::VERSION
	s.platform      = Gem::Platform::RUBY
	s.license       = "MIT"
	s.authors       = ['Scott Rabin']
	s.email         = ['showrobot@scottrabin.com']
	s.homepage      = ''
	s.summary       = 'An easy-to-use command line tool for sorting, classifying, and organizing your media'
	s.description   = ''

	s.required_rubygems_version = '>= 1.3.6'
	s.rubyforge_project         = 'showrobot'

	s.add_development_dependency 'rspec'

	s.files         = Dir.glob("{bin,lib}/**/*") + %w(LICENSE README.md)
	s.executables   = ['showrobot']
	s.require_path  = 'lib'
end
