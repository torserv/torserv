# TorServ Firejail profile v1.0 (tested on Firejail 0.9.72)
# Future versions of Firejail may fail to start with this profile

# torserv.profile
private
private-dev
read-only tor
read-only public
noexec /tmp
caps.drop all
seccomp
nonewprivs
