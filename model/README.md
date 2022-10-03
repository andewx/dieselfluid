# SPH FLuid Model

## Overview

dslfluid/model package provides interfaces and functional patterns for working
with fluid libraries. Also provides a special interface to the GPU

### Notes

- Density field kernel methods will utilize a new correction method for Smoothed Particle fields
known as the Spatial Corrective Variance Procedure. Step 1: Calculated Per Particle Guassian Neighbor Models. Guassian Mean + Variance should determine stability of the density field estimator. For the corrective procedure we estimate the density valence of each position and
initialize a ML error correction algorithm based on MLP which generates target variance augmentation of the estimated neighborhood means.

- For the Smoothed Particle Method the EOS function and weighted calculations for
the density field estimator will be based on a sampled neighbor space for x[N]
parameter particles. The distribution of neighbor distances will be sampled and
incorporated into a guassian distribution with variance o and mean u. The Density field
distributions will then estimate a corrective variance bias for adjusting and smoothing
high density fields. These will be corrective variance distributions and accumulate throughout
the initialization and correction procedure for a Spatial Corrective Variance Procedure.
