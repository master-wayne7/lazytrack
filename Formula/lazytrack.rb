class Lazytrack < Formula
  desc "A fun CLI-based time/habit tracker"
  homepage "https://github.com/master-wayne7/lazytrack"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/master-wayne7/lazytrack/releases/download/v1.0.0/lazytrack_Darwin_arm64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    else
      url "https://github.com/master-wayne7/lazytrack/releases/download/v1.0.0/lazytrack_Darwin_x86_64.tar.gz"
      sha256 "E0EECC7010E5406AE40E3DF3E9D09E05627FC8FEE288505311C20289EAFF2D06"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/master-wayne7/lazytrack/releases/download/v1.0.0/lazytrack_Linux_arm64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    else
      url "https://github.com/master-wayne7/lazytrack/releases/download/v1.0.0/lazytrack_Linux_x86_64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end
  end

  def install
    bin.install "lazytrack"
  end

  test do
    system "#{bin}/lazytrack", "--help"
  end
end 