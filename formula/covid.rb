# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Covid < Formula
  desc "This tool just downloads the RKI covid raw data from link:https://media.githubusercontent.com/media/robert-koch-institut/SARS-CoV-2_Infektionen_in_Deutschland/master/Aktuell_Deutschland_SarsCov2_Infektionen.csv[] and shows them on the command line."
  homepage ""
  version "0.2"
  license "Apache 2.0 License"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/steffakasid/covid/releases/download/v0.2/covid_0.2_Darwin_x86_64.tar.gz"
      sha256 "c81e246587d41ef19d60df802e85e9c7e88c7548780fcd4cf24a23f555e15ba8"

      def install
        bin.install "covid"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/steffakasid/covid/releases/download/v0.2/covid_0.2_Darwin_arm64.tar.gz"
      sha256 "952f6a15acc4f16ddebefb8bed1992d604a2ba205a46a284f699f6df2b721940"

      def install
        bin.install "covid"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/steffakasid/covid/releases/download/v0.2/covid_0.2_Linux_x86_64.tar.gz"
      sha256 "f59413a21f8e4db9434032655abfa972341c8c4c5bbfcfde989b75e023b77ba9"

      def install
        bin.install "covid"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/steffakasid/covid/releases/download/v0.2/covid_0.2_Linux_arm64.tar.gz"
      sha256 "44ef7edc251645857b6b31bc140ce83ddf3d695d717542b0650751cc35da98c5"

      def install
        bin.install "covid"
      end
    end
  end
end
