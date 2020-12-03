use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

fn count_tree(pattern: &Vec<char>, x: usize) -> i32 {
  let l = pattern.len();
  if pattern[x%l] == '#' {
    return 1;
  }
  return 0;
}

fn main() {
    if let Ok(lines) = read_lines("./day3_input.txt") {
      let mut x: usize = 0;
      let mut sum: i32 = 0;
      for line in lines {
          if let Ok(l) = line {
            let char_vec: Vec<char> = l.chars().collect();
            sum += count_tree(&char_vec, x);
            x += 3;
          }
      }
      println!("Num Trees: {}", sum);
    }
}

// The output is wrapped in a Result to allow matching on errors
// Returns an Iterator to the Reader of the lines of the file.
fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}