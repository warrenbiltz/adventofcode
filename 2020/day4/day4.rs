use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;
use std::collections::HashMap;
extern crate regex;
use regex::Regex;

fn parse_value(passport: &mut HashMap<String, String>, value_pair: &str) {
    let vp: Vec<&str> = value_pair.split(':').collect();
    if vp.len() == 2 {
      passport.insert(vp[0].to_string(), vp[1].to_string());
    }
}

fn parse_line(mut passport: &mut HashMap<String, String>, line: &str) {
  let value_pairs = line.trim().split_whitespace();
  for vp in value_pairs {
    parse_value(&mut passport, vp);
  }
}

fn validate_yr(passport: &HashMap<String, String>, key: &str, min: i32, max: i32) -> bool {
  match passport.get(key) {
    Some(value) => {
      let re = Regex::new(r"^\d{4}$").unwrap();
      if re.is_match(value) {
        let year = value.parse::<i32>().unwrap();
        if year >= min && year <= max {
          return true
        }
      }
      return false
    },
    None => return false
  }
}

fn validate_hgt(passport: &HashMap<String, String>) -> bool {
  match passport.get("hgt") {
    Some(value) => {
      let re = Regex::new(r"^(\d+)(cm|in)$").unwrap();
      if re.is_match(value) {
        let caps = re.captures(value).unwrap();
        let num = caps.get(1).unwrap().as_str().parse::<i32>().unwrap();
        let unit = caps.get(2).unwrap().as_str();
        if unit == "in" && num >= 59 && num <= 76 {
          return true
        }
        if unit == "cm" && num >= 150 && num <= 193 {
          return true
        }
      }
      return false
    },
    None => return false
  }
}

fn validate_regex(passport: &HashMap<String, String>, key: &str, exp: &str) -> bool {
  match passport.get(key) {
    Some(value) => {
      let re = Regex::new(exp).unwrap();
      return re.is_match(value)
    },
    None => return false
  }
}

fn validate_passport2(passport: &HashMap<String, String>) -> i32 {
  if !validate_yr(&passport, "byr", 1920, 2002) {
    return 0;
  }
  if !validate_yr(&passport, "iyr", 2010, 2020) {
    return 0;
  }
  if !validate_yr(&passport, "eyr", 2020, 2030) {
    return 0;
  }
  if !validate_hgt(&passport) {
    return 0;
  }
  if !validate_regex(&passport, "hcl", r"^#[a-f0-9]{6}$") {
    return 0;
  }
  if !validate_regex(&passport, "ecl", r"^(amb)|(blu)|(brn)|(gry)|(grn)|(hzl)|(oth)$") {
    return 0;
  }
  if !validate_regex(&passport, "pid", r"^\d{9}$") {
    return 0;
  }
  return 1;
}

fn validate_passport1(passport: &HashMap<String, String>) -> i32 {
  if !passport.contains_key("byr") {
    return 0;
  }
  if !passport.contains_key("iyr") {
    return 0;
  }
  if !passport.contains_key("eyr") {
    return 0;
  }
  if !passport.contains_key("hgt") {
    return 0;
  }
  if !passport.contains_key("hcl") {
    return 0;
  }
  if !passport.contains_key("ecl") {
    return 0;
  }
  if !passport.contains_key("pid") {
    return 0;
  }
  return 1;
}

fn main() {
    if let Ok(lines) = read_lines("./day4_input.txt") {
      let mut passport = HashMap::<String, String>::new();
      let mut valid1 = 0_i32;
      let mut valid2 = 0_i32;
      for line in lines {
          if let Ok(l) = line {
            if !l.trim().is_empty() {
              parse_line(&mut passport, &l);
            }
            else {
              valid1 += validate_passport1(&passport);
              valid2 += validate_passport2(&passport);
              passport.clear();
            }
          }
      }
      valid1 += validate_passport1(&passport);
      valid2 += validate_passport2(&passport);
      println!("Valid Passports: 1 = {} | 2 = {}", valid1, valid2);
    }
}

// The output is wrapped in a Result to allow matching on errors
// Returns an Iterator to the Reader of the lines of the file.
fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}