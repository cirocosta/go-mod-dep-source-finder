local repositories = std.extVar('repositories');

local repositoriesResources = [{
  name: repository.name + '-repo',
  type: 'git',
  source: { uri: repository.uri },
} for repository in repositories];

local repositoriesJobs = [{
  name: repository.name + '-deps',
  public: true,
  serial: true,
  plan: [
    {
      aggregate: [
        {
          get: 'repository',
          passed: ['build'],
          trigger: true,
        },
        {
          get: 'rc-image',
          passed: ['build'],
          trigger: true,
        },
        {
          get: repository.name + '-repo',
          trigger: true,
        },
      ],
    },
    {
      task: 'gather-dependencies',
      input_mapping: { repo: repository.name + '-repo' },
      file: 'repository/ci/tasks/gather-dependencies-list.yml',
    },
    {
      task: 'run',
      image: 'rc-image',
      file: 'repository/ci/tasks/resolve-deps.yml',
    },
  ],
} for repository in repositories];

{
  resources: repositoriesResources + [
    {
      name: 'repository',
      type: 'git',
      source: { uri: 'https://((github-token))@github.com/cirocosta/go-mod-dep-source-finder' },
    },
    {
      name: 'alpine',
      type: 'registry-image',
      source: { repository: 'alpine' },
    },
    {
      name: 'golang-alpine',
      type: 'registry-image',
      source: { repository: 'golang', tag: 'alpine' },
    },
    {
      name: 'rc-image',
      type: 'registry-image',
      source: {
        repository: 'cirocosta/go-mod-dep-source-finder',
        username: '((docker-user))',
        password: '((docker-password))',
        tag: 'rc',
      },
    },
    {
      name: 'final-image',
      type: 'registry-image',
      source: {
        repository: 'cirocosta/go-mod-dep-source-finder',
        username: '((docker-user))',
        password: '((docker-password))',
      },
    },
  ],
  jobs: repositoriesJobs + [
    {
      name: 'build',
      public: true,
      serial: true,
      plan: [
        {
          aggregate: [
            { get: 'alpine', trigger: true },
            { get: 'golang-alpine', trigger: true },
            { get: 'repository', trigger: true },
          ],
        },
        {
          task: 'build-image',
          privileged: true,
          file: 'repository/ci/tasks/build-image.yml',
        },
        {
          put: 'rc-image',
          inputs: ['image'],
          get_params: { format: 'oci' },
          params: { image: 'image/image.tar' },
        },
      ],
    },
    {
      name: 'publish-image',
      public: true,
      serial: true,
      plan: [
        {
          get: 'rc-image',
          params: { format: 'oci' },
          passed: [repository.name + '-deps' for repository in repositories],
          trigger: false,
        },
        {
          put: 'final-image',
          params: { image: 'rc-image/image.tar' },
        },
      ],
    },
  ],
}
