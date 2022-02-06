# -*- mode: ruby -*-
# vi: set ft=ruby :

synced_folder  = ENV[     'SYNCED_FOLDER'      ]  || "/home/vagrant/#{File.basename(Dir.pwd)}"
memory         = ENV[           'MEMORY'       ]  || 4096
cpus           = ENV[           'CPUS'         ]  || 4
vm_name        = ENV[           'VM_NAME'      ]  || File.basename(Dir.pwd)
forwarded_ports= []
provisioners   = [
  "node",
  "python",
  "ripgrep",
  "docker",
  "starship",
  "rust-core-utils",
  "rust-toolchain",
  "spacevim",
]
$cleanup_script = <<-SCRIPT
apt-get autoremove -yqq --purge > /dev/null 2>&1
apt-get autoclean -yqq > /dev/null 2>&1
apt-get clean -qq > /dev/null 2>&1
rm -rf /var/lib/apt/lists/*
SCRIPT
INSTALLER_SCRIPTS_BASE      = "https://raw.githubusercontent.com/da-moon/provisioner-scripts/master/bash/installer"
UTIL_SCRIPTS_BASE           = "https://raw.githubusercontent.com/da-moon/provisioner-scripts/master/bash/util"
Vagrant.configure("2") do |config|
  config.vm.define "#{vm_name}"
  config.vm.hostname = "#{vm_name}"
  config.vm.synced_folder ".","#{synced_folder}",auto_correct:true, owner: "vagrant",group: "vagrant",disabled:true
  config.vm.provider "virtualbox" do |vb, override|
    override.vm.box="generic/debian10"
    vb.memory = "#{memory}"
    vb.cpus   = "#{cpus}"
    # => enable nested virtualization
    vb.customize ["modifyvm",:id,"--nested-hw-virt", "on"]
    override.vm.synced_folder ".", "#{synced_folder}",disabled: false,
      auto_correct:true, owner: "vagrant",group: "vagrant",type: "virtualbox"
  end if Vagrant.has_plugin?('vagrant-vbguest')
  config.vm.provider "hyperv" do |h,override|
    override.vm.box="generic/debian10"
    h.enable_virtualization_extensions = true
    h.linked_clone = true
    h.cpus   = "#{cpus}"
    h.memory = "#{memory}"
    h.maxmemory = "#{memory}"
    override.vm.network "public_network"
    override.vm.synced_folder ".", "#{synced_folder}",disabled: false,auto_correct:true, type: "smb",
    owner: "vagrant",group: "vagrant"
  end
  config.vm.provider "libvirt" do |libvirt,override|
    override.vm.box="generic/debian10"
    libvirt.memory = "#{memory}"
    libvirt.cpus = "#{cpus}"
    libvirt.nested = true
    libvirt.cpu_mode = "host-passthrough"
    libvirt.driver = "kvm"
    override.vm.synced_folder ".", "#{synced_folder}",
      disabled: false,auto_correct:true, owner: "1000", group: "1000",
      type: "9p",accessmode: "passthrough"
    override.vm.provision "shell",
      privileged:true,
      name:"p9-kernel-support",
      inline: <<-SCRIPT
      [ ! -L /usr/local/bin/modprobe ] && sudo ln -s /sbin/modprobe /usr/local/bin/modprobe
      SCRIPT
  end if Vagrant.has_plugin?('vagrant-libvirt')
  # [ NOTE ] => vagrant 2.2.16 has a bug in which makes it impossible to provision
  # google boxes with ssh , either use version >=2.2.17 or version <=2.2.16
  config.vm.provider "google" do |google,override|
    google.name                       = "#{vm_name}"
    google.disk_type                  = "pd-ssd"
    google.disk_size                  = ENV["GCLOUD_DISK_SIZE"] || "50"
    google.google_json_key_location   = ENV['GOOGLE_APPLICATION_CREDENTIALS'] || File.expand_path("~/#{vm_name}_gcloud.json").chomp()
    google.google_project_id          = ENV['GOOGLE_PROJECT_ID'] || %x[gcloud config get-value core/project].chomp()
    google.zone                       = ENV["CLOUDSDK_COMPUTE_ZONE"] || %x[gcloud config get-value compute/zone].chomp()
    google.machine_type               = ENV["GCLOUD_MACHINE_TYPE"] || "n1-standard-8"
    google.image                      = ENV["GCLOUD_IMAGE"] || %x[gcloud compute images list --format='value(NAME)' --filter='name ~ debian AND family ~ debian-10'].chomp()
    google.image_project_id           = ENV["GCLOUD_IMAGE_PROJECT_ID"] || %x[gcloud compute images list --format='value(PROJECT)' --filter='name ~ debian AND family ~ debian-10'].chomp()
    google.metadata                   = {'ssh-keys' => "#{ENV['USERNAME'] || %x[whoami].chomp()}:#{File.read(File.expand_path("~/.ssh/#{vm_name}.pub").chomp()).chomp()}"}
    override.vm.box                   = "google/gce"
    override.ssh.username             = ENV['USERNAME'] || %x[whoami].chomp()
    override.ssh.private_key_path     = File.expand_path("~/.ssh/#{vm_name}")

    override.vm.synced_folder ".", "/workspace", type: 'rsync',
      rsync__args: ["--verbose", "--archive", "-z"],
      owner: ENV['USERNAME'] || %x[whoami].chomp(),
      group: ENV['USERNAME'] || %x[whoami].chomp(),
      sync__exclude: [
        '.vagrant',
      ]
  end if Vagrant.has_plugin?('vagrant-google')
  forwarded_ports.each do |port|
    config.vm.network "forwarded_port",
      guest: port,
      host: port,
      auto_correct: true
  end
  config.vm.provision "shell",privileged:true,name:"cleanup", inline: $cleanup_script
  config.vm.provision "shell",
    privileged:false,
    name:"init",
    path: "#{INSTALLER_SCRIPTS_BASE}/init"
  # [ NOTE ] => provisioning
  provisioners.each do |provisioner|
    config.vm.provision "shell",
      privileged:false,
      name:"#{provisioner}",
      path: "#{INSTALLER_SCRIPTS_BASE}/#{provisioner}"
  end
  config.vm.provision "shell",
      privileged:false,
      name:"extra-tools",
      inline: <<-SCRIPT
			set -xeu
      sudo apt-get install -y upx cmake libssl-dev fzf libgconf-2-4  ;
      for i in {1..5}; do wget -O \
        /tmp/vsls-reqs \
        https://aka.ms/vsls-linux-prereq-script && break || sleep 15; done ;
      sudo bash /tmp/vsls-reqs ;
      rm -f /tmp/vsls-req ;
      sudo snap install diagon ;
      sudo python3 -m pip install asciinema yq pre-commit
      sudo yarn global add --silent --prefix /usr/local \
        @commitlint/cli \
        @commitlint/config-conventional \
        remark \
        remark-cli \
        remark-stringify \
        remark-frontmatter \
        wcwidth \
        prettier \
        bash-language-server \
        dockerfile-language-server-nodejs
      rustup default stable
      cargo install -j`nproc` convco ;
    SCRIPT
  config.trigger.after [:provision] do |t|
    t.info = "cleaning up after provisioning"
    t.run_remote = {inline: $cleanup_script }
  end
end
