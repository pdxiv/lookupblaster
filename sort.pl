#!/usr/bin/perl
use strict;
use warnings;
our $VERSION = '0.0.1';

my @line;
while (<>) {
    push @line, $_;
}

my @ips;
my %data;
foreach (@line) {
    my ( $key, $value ) = split /\s+/msx;
    $data{$key} = $value;
    push @ips, $key;
}

@ips = map { sprintf '%d.%d.%d.%d', split /\./msx } sort map { sprintf '%03d.%03d.%03d.%03d', split /\./msx } @ips;

foreach (@ips) {
    print "$_\t$data{$_}\n";
}
