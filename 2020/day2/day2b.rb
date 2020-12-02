def is_pw_valid?(line)
  match = /(\d+)-(\d+)\s+(\S):\s+(\S+)/.match(line)
  return unless match
  char1 = match[4][match[1].to_i - 1]
  char2 = match[4][match[2].to_i - 1]
  return char1 != char2 && (char1 == match[3] || char2 == match[3])
end

File.open("./day2_input.txt", "r") do |f|
  num_valid = f.select { |line| is_pw_valid? line }.size
  puts "VALID: #{num_valid}"
end
