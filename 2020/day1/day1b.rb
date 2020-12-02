require 'set'

cache1 = Set.new
cache2 = {}

File.open("./day1a_input.txt", "r") do |f|
  f.each_line do |line|
    num = line.to_i
    diff = 2020 - num
    if cache2.has_key? diff
      puts "FOUND #{num} * #{cache2[diff][1]} * #{cache2[diff][2]} = #{num * cache2[diff][0]}"
      return num * cache2[diff][0]
    else 
      cache1.each do |x|
        cache2[x+num] = [x * num, x, num]
      end
      cache1.add(num)
    end
  end
end