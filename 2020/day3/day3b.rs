use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

fn count_trees(map: &Vec<Vec<char>>, slope_x: usize, slope_y: usize) -> i64 {
  let mut x = 0_usize;
  let mut y = 0_usize;
  let mut sum = 0_i64;

  while y < map.len() {
    let l = map[y].len();
    if map[y][x%l] == '#' {
      sum += 1;
    }
    x += slope_x;
    y += slope_y;
  }
  println!("count_trees: {}/{} = {}", slope_x, slope_y, sum);
  return sum;
}

fn main() {
    if let Ok(lines) = read_lines("./day3_input.txt") {
      let mut map = Vec::<Vec<char>>::new();
      for line in lines {
          if let Ok(l) = line {
            map.push(l.chars().collect::<Vec<char>>());
          }
      }
      let mut prod: i64 = 1;
      let slopes = vec![ vec![1,1], vec![3,1], vec![5,1], vec![7,1], vec![1,2]];
      for slope in slopes {
        prod *= count_trees(&map, slope[0], slope[1]);
      }
      println!("Prod Trees: {}", prod);
    }
}

// The output is wrapped in a Result to allow matching on errors
// Returns an Iterator to the Reader of the lines of the file.
fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}