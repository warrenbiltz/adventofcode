require 'set'

cache = Set.new

File.open("./day1a_input.txt", "r") do |f|
  f.each_line do |line|
    num = line.to_i
    diff = 2020 - num
    if cache.include? diff
      puts "FOUND #{num + diff}: #{num} * #{diff} = #{num * diff}"
      return num * diff
    end
    cache.add(num)
  end
end