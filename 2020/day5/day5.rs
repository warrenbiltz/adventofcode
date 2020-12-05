use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

fn decode(line: &str) -> (i32, i32, i32) {
  let seat = line.chars().collect::<Vec<char>>();
  let mut row = 127_i32;
  let mut col = 7_i32;
  let mut partition = 64;

  let mut i = 0;
  while i < 7 {
      if seat[i] == 'F' {
        row -= partition;
      }
      partition /= 2;
      i += 1;
  }
  partition = 4;
  while i < 10 {
      if seat[i] == 'L' {
        col -= partition;
      }
      partition /= 2;
      i += 1;
  }

  return (row, col, 8 * row + col);
}

fn main() {
  if let Ok(lines) = read_lines("./day5_input.txt") {
    let mut min_seat_id = (128, 8, 1024);
    let mut max_seat_id = (0, 0, 0);
    let mut id_sum = 0;
    for line in lines {
      if let Ok(l) = line {
        let seat_id = decode(&l);
        id_sum += seat_id.2;
        if seat_id.2 > max_seat_id.2 {
          max_seat_id = seat_id;
        }
        if seat_id.2 < min_seat_id.2 {
          min_seat_id = seat_id;
        }
      }
    }
    let m = min_seat_id.2;
    let n = max_seat_id.2;
    let expected_id_sum = (m + n) * (n - m + 1) / 2;
    println!("min seat: {:?}", min_seat_id);
    println!("max seat: {:?}", max_seat_id);
    println!("missing_seat: {}", expected_id_sum - id_sum);
  }
}

// The output is wrapped in a Result to allow matching on errors
// Returns an Iterator to the Reader of the lines of the file.
fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}