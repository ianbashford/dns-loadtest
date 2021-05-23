FROM archlinux:base-devel

RUN pacman -Sy --noconfirm \
	curl \
	tar \
	vim \
	dnsperf

# Install yay - https://github.com/Jguer/yay
ENV yay_version=9.4.2
ENV yay_folder=yay_${yay_version}_x86_64

RUN cd /tmp && \
    curl -L https://github.com/Jguer/yay/releases/download/v${yay_version}/${yay_folder}.tar.gz | tar zx && \
    install -Dm755 ${yay_folder}/yay /usr/bin/yay && \
    install -Dm644 ${yay_folder}/yay.8 /usr/share/man/man8/yay.8 


# Create a new user with no password
ENV USERNAME developer
RUN set -ex && \
	useradd --create-home --key MAIL_DIR=/dev/null -G wheel --shell /bin/bash $USERNAME && \
	passwd -d $USERNAME && \
	echo "%wheel ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

USER developer

CMD [ "/bin/bash" ]
