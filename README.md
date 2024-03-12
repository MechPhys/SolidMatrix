# SolidMatrix: Solid particles in a matrix

> "Nature is, above all else, a meticulous accountant"

SolidMatrix is a new approach that is based on discrete particles in a discrete matrix or lattice.  These discrete particles are much better for keeping the accounting in order (i.e., strict conservation of energy, charge, particle numbers, etc), but how does a hard little discrete dot move in a seemingly continuous fashion in arbitrary directions?  Interestingly, stochastic motion, where there is a certain probability of taking a step at each point in time, provides a way of smoothing the motion of these discrete particles.  This provides a necessary origin story for the randomness at the heart of quantum physics, and does a great job of keeping the accounting tight.

In particular, the discrete localization of particles to one cell of a matrix allows for a complete energy-conserving accounting of particle interactions, in a way that eludes any distributed particle system, where it is very difficult and non-local to "gather all the far-flung bits of a particle" and alter them in some systematic way necessary to conserve energy during an interaction.

The resulting physical model is like the deBroglie-Bohm pilot wave theory: the particles interact with a sea of surrounding fields that impart the wave-like properties of QM.  Antonio Sciarretta has developed models based on this principle, in a series of recent papers, showing how a particular formulation of the surrounding field interactions can give rise to otherwise puzzling features of quantum physics.  In particular, these field interactions introduce a novel framework for understanding interference effects, in terms of the cumulative traces left by other particles.  This twist may well provide the critical missing ingredient that finally makes the quantum world make sense: it is well-established that quantum effects manifest only over many repeated samples from "identically prepared" particles, and yet everyone interprets them in terms of physical mechanisms operating only on a single trial, associated with a single isolated particle.  This fundamental disconnect has great potential for resolving paradoxes.

This approach also converges with the research on the zero-point field, as developed by Santos, Marshall and others, which is synonymous with the virtual particles captured in the path-integral formulation of Feynmann and his famous diagrams.  The field is very much "alive", teaming with random "quantum foam", and this fabric of "empty" space likely plays an essential role in explaining various puzzling features of QM.  In the stochastic electrodynamics (SED) approach, the field is the exclusive source of stochastic behavior.  However, the isotrophic motion of discrete particles requires intrinsic stochasticity, thus freeing the field to play a less strongly-constrained role.

In short, we have several different models of the quantum field: the original S field formulation from Bohm, the path integral of Feynmann (which is divergent in general), the second quantized Fourier space formulation of quantum field theory, the SED / ZPF field, and Sciarretta's specific "trace" field.  Between all of these models, some kind of consistent framework should be derivable.

## Virtual particles, the discrete lattice, and probability waves

Virtual particles are an essential feature of QED / quantum field theory, and yet their "ontological" status is clearly somewhat confusing: they aren't the "real" particles that we observe, and yet their fleeting existence is necessary for the theory to work, so in some sense they must be just as real as the "real" particles.

The discrete particle lattice framework provides a potential resolution to this conundrum.  If any given "real" particle can potentially occupy any given cell in a discrete lattice, then there must effectively be a "slot" reserved for each such particle in each cell.  These empty slots could provide an appealing basis for virtual particles, and the propagation and interactions of particles in the matrix.

In particular, a simple schema is that the probability waves associated with the standard interpretation of QM reflect a rippling propagation of probability factors across virtual particle slots in the matrix, with a real particle having a "1" special status.  Each possible jump to a neighboring cell involves a full transition matrix dependent upon the total energy (mass + kinetic) of the source: if the source is sufficiently energetic, it has some probability of activating a different combination of real particles as it makes the leap, accounting for the splitting tracks observed in particle accelerator experiments.  Perhaps some of the "trace" in the matrix represents residual bits of this probability field propagating out and being left behind as real particles move around.

It is essential that these probability computations are all propagated in terms of _amplitudes_, not the probability values themselves, which are obtained by the product with the complex conjugate ("squaring").

The force-carrying fields / particles need to be incorporated into this scheme.  For "photons" in the EM field, they do not mutually interact and would seem to require a continuous-valued representation of the field to represent the varying wave intensity or amplitude.  It is possible that the EM field is essentially classical as a real-valued field, but because its energy contributes to the particle probability field, the discrete dynamics of that field essentially creates photon-like particle effects whenever it shows itself on something actually measurable.

The frequency of the EM wave determines the energy, and this is well represented in the standard wave equation Hamiltonian, so that is presumably the property that matters in influencing the particle probability field.  However, it is unclear if this would be locally measurable in any meaningful way, and accounting-wise properly subtracted from the ongoing field, so perhaps it is actually necessary to treat photons as real particles too?

## The static programmability constraint

The overall goal of this approach can be summarized in terms of the properties of a computer program that implements it.  The best case scenario is that all the seemingly arbitrary and puzzling aspects of quantum physics fall out naturally from this perspective, providing a unifying overall framework for _how_ Nature could actually run physics in a fully autonomous, self-contained manner.  By contrast, existing physics frameworks require case-specific and incomplete mathematical approaches, and only generate approximate solutions to particular physical situations.

If such a framework could be developed, it might also go a long way toward understanding _why_ physics is the way it is. On the other hand, the ultimate big-picture question of why it is all here in the first place will not likely succumb to any simple explanation.

In particular, an autonomous implementation of physics must:

* Use a comprehensive "static" state representation that can contain all the relevant degrees of freedom, _without requiring dynamic memory allocation_ at random points.  We can't really have Nature dipping into some kind of memory allocation pool and dishing out bits of relevant state here and there as needed.  Everything needs to have its place, and physics is just the updating of this state matrix.

* Interactions are fully local (nearest neighbor on the grid), as in a cellular automaton, because again the alternative would require arbitrary subroutines of processing to be performed in different locations at different points in time: how would all this get coordinated and allocated?  Autonomous physics means _no daemons_ -- everything "just runs" automatically, following the same rules over and over again.

In this framing, the goal is to come up with the probability-amplitude based equations for how each particle state moves through the matrix over time.  If this can be done for low-energy electrons in the EM field, as a starting point, then presumably it could be generalized to more types of particles etc.  The heavier muon and tau leptons, and the corresponding neutrinos, are potentially all reducible to one underlying particle type, which is variously "decorated" or parameterized with respect to mass and charge, and their interactions with the weak force.  Working through all of that, along with the basis for spin, would be a logical next step.

## Measurement and wavefunction collapse

The specific example of the Born rule and the mysterious collapse of the wave function at the point of _measurement_, which remains one of the most important unsolved conceptual problems at the heart of quantum physics, is particularly illustrative for the defining features of this approach.

Under the classical Copenhagen interpretation, Schroedinger etc wave functions define the evolution of a _probability field_ for where a particle might be found at some point in time.  But at some imprecisely-defined moment of "measurement", this probability field collapses down to a single definite point, where the particle was actually observed.  But how could this spatially distributed field, which could spread out without bound and in principle cover the entire universe, _instantaneously collapse_ to a single point?  It just makes no physical sense, even while it makes practical mathematical sense in terms of predicting outcomes of experiments.

From a programming perspective, it is completely untenable.  Imagine the for loop that would be necessary to track down all the far-flung bits of the probability field, and the challenges of dealing with very tiny values (is there a cutoff?) -- it would require a massive computational load for something that is presumably happening everywhere, all the time, in parallel throughout the universe.

The obvious alternative is that the probability wave describes the behavior of a local, presumably stochastic process, _where the particle remains definitively localized at all times_.  This framework avoids any need for wave function collapse, and could be directly implemented in a computer simulation.  The deBroglie-Bohm _pilot wave_ framework advanced this idea, where a particle is guided in its trajectory by some kind of wave function, but retains a definite location at all times.

However, to explain the non-local _entanglement_ phenomena in QM, this pilot wave function must have a dimensionality that grows exponentially as a function of the number of interacting particles.  Because there is no in-principle limit to the number of such particles, in principle the entire Universe must enter into this exponentially huge wave function, which is then inexpressably huge.  Bohm and Hiley refer to this as the "undivided universe".

Critically, it is possible that the trace mechanisms, along with other potential propagation of interaction terms in the updating of probability factors in the particle matrix theory, could account for the interaction terms that the exponentially-expanding state function in the pilot wave framework requires.  In this context, it is essential to imagine a churning sea of virtual quantum foam interacting with neighbors at light speed -- the idea that such a thing could be represented in summary form by a low-dimensional linear set of probability amplitude values seems implausible at best.

Thus, the critical refrain here is that there are all manner of simple mathematical models in physics that are almost certainly not viable _physical models_ for how physics actually unfolds autonomously over time.  Here, we seek the small (unique?) subset of models that makes good physical sense in terms of the above programmability constraints.

## One particle stochastic equations of motion

## Outstanding challenges

Sciarretta's specific formulation has some strange features, including that the interaction depends on the specific lifetime of the particle.  However, in practice the field is initialized to an equilibrium distribution, based on the configuration.

Thus, in the end, the field "senses" the configuration and imparts corresponding forces onto the particle, which is the same qualitative account as the pilot wave framework.  Furthermore, Sciarretta's model has an appealing conserving exchange process in the particle's interaction with the field.

Some other challenges for the field model include accounting for the space-time distortion effects of special relativity, which are so elegantly captured in the wave dynamics.

Although hard discrete particles make sense for fermions, which obey strict conservation principles and the Pauli exclusion principle, bosons like the troubling photon have very different properties.  They have no mass or charge, accumulate in graded amounts with fully linear superposition, always move at the speed of light, and the relevant wavelengths can be astronomically large compared to the Compton wavelength of even the lightest particles.

![Stochastic Origin of Quantum Momentum / Frequency Relationship](figs/fig_asmom5_0_autoc.png?raw=true "Stochastic Origin of Quantum Momentum / Frequency Relationship.  The momentum on the left is 0.5c while on the right is 0. The distribution of position is on the vertical axis, while time is on the horizontal axis, with each point centered at the origin in the center (i.e., the temporal autocorrelation function).  The variance on the left is half of that on the right.")

Nelson (1966) shows how brownian stochastic motion gives rise to the Schrodinger equation: a central intuition is that slow drift produces a wide cloud of space where particle could be, corresponding to a long wavelength in the probability cloud that the Schrodinger wave function describes.  However, when the particle has high momentum, it moves more deterministically in the given direction, resulting in a narrower range of variance around the particle's mean trajectory, resulting in a narrower effective wavelength along the direction of motion.  This is illustrated in the above figure.

## Differences Between Photons and Real Particles

Photons differ in so many ways from electrons and other particles, it is hard to imagine that both could be described by the same underlying mechanistic model.  Rest mass, charge, conserved number, fractional spin, etc.  In short, it is difficult to see how photons could be fit within the same paradigm.  The simplest approach at this point is to implement EM using classical Maxwell's equations.

# References

See https://www.zotero.org/groups/2525742/mechphys/library for an extensive library of relevant papers.

* O’Reilly, R. C. (2011). Surely You Must All be Joking: An Outsider’s Critique of Quantum Physics. ArXiv:1109.0880]. http://arxiv.org/abs/1109.0880

* Santos E. (2015). Towards a realistic interpretation of quantum mechanics providing a model of the physical world, Foundations of Science 20(4), 357–386. https://doi.org/10.1007/s10699-014-9366-y

* Sciarretta, A. (2018). A Local-Realistic Model of Quantum Mechanics Based on a Discrete Spacetime. Foundations of Physics, 48(1), 60–91. https://doi.org/10.1007/s10701-017-0129-9

* Sciarretta, A. (2018). A Local-Realistic Model of Quantum Mechanics Based on a Discrete Spacetime (Extended version). Foundations of Physics, 48(1), 60–91. http://arxiv.org/abs/1712.03227

* Sciarretta, A. (2021). A local-realistic quantum mechanical model of spin and spin entanglement. International Journal of Quantum Information, 19(01), 2150006. https://doi.org/10.1142/S0219749921500064


