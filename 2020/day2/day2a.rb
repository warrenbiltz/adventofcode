def is_pw_valid?(line)
  match = /(\d+)-(\d+)\s+(\S):\s+(\S+)/.match(line)
  return unless match
  count = match[4].count(match[3])
  return count >= match[1].to_i && count <= match[2].to_i
end

File.open("./day2_input.txt", "r") do |f|
  num_valid = f.select { |line| is_pw_valid? line }.size
  puts "VALID: #{num_valid}"
end
