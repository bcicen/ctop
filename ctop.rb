class Ctop < Formula
  desc "Top-like interface for container metrics"
  homepage "https://bcicen.github.io/ctop/"
  url "https://github.com/bcicen/ctop/releases/download/v0.4.1/ctop-0.4.1-darwin-amd64"
  sha256 "884555d6303652ba7892ee90faa6527bacbfb73b84c1006edd8916bafb5de22f"

  bottle :unneeded

  def install
    mv "ctop-0.4.1-darwin-amd64", "ctop"
    bin.install "ctop"
  end

  test do
    system "#{bin}/ctop", "version"
  end
end
