# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v4.6.0] - 2021-09-30

### Added

- test: Add tests for GSP-749 unify path behavior (#59)
- test: Add IoCallback for read and write tests (#61)
- test: Add integration tests for MultipartHTTPSigner (#62)
- test: Add read with offset and size tests (#63)

### Fixed

- fix: Replace backslash with slash for path behavior test (#60)

## [v4.5.0] - 2021-09-23

### Added

- storage: Implement GSP-751 add Write Empty File Behavior integration tests (#56)

## [v4.4.0] - 2021-09-03

### Added

- storage: Implement StorageHTTPSigner test (#49)

## [v4.3.0] - 2021-08-02

### Added

- link: Implement GSP-86 add linker integration tests (#40)

### Changed

- append: Add test for appending to an exists file (#44)
- copy: Split TestCopier (#42)
- move: Split TestMover (#42)

### Fixed

- link: Fix target check failed for linker (#47)

## [v4.2.0] - 2021-07-14

### Added

- storage: Implement Copier and Mover tests (#27)

### Changed

- storage: Minor refactor (#31)
- storage: Implement GSP-134 Write Behavior Consistency (#32)
- storage: Update tests for List (#33)

### Fixed

- move: Use errors.Is to assert error's type (#28)
- append: Remove content-length check for CreateAppend (#34)

### Upgraded

- build(deps): bump github.com/google/uuid from 1.2.0 to 1.3.0 (#35)

## [v4.1.1] - 2021-06-25

### Changed

- dir: Add path check for direr (#23)

### Fixed

- multipart: Fix CompletePart not passing (#22)

## [v4.1.0] - 2021-06-10

### Added

- storage: Implement Direr tests (#20)

## [v4.0.0] - 2021-05-24

### Changed

- storage: Add CommitAppend (#15)
- storage: Implement GSP-46 Idempotent storager delete operation (#16)
- storage: Implement GSP-62 WriteMultipart returns Part (#17)
- *: Implement GSP-73 Organization Rename (#18)

## v3.0.0 - 2021-04-20

### Added

- Implement integration test for storager
- Implement integration test for multiparter
- Implement integration test for appender

[v4.6.0]: https://github.com/beyondstorage/go-integration-test/compare/v4.5.0...v4.6.0
[v4.5.0]: https://github.com/beyondstorage/go-integration-test/compare/v4.4.0...v4.5.0
[v4.4.0]: https://github.com/beyondstorage/go-integration-test/compare/v4.3.0...v4.4.0
[v4.3.0]: https://github.com/beyondstorage/go-integration-test/compare/v4.2.0...v4.3.0
[v4.2.0]: https://github.com/beyondstorage/go-integration-test/compare/v4.1.1...v4.2.0
[v4.1.1]: https://github.com/beyondstorage/go-integration-test/compare/v4.1.0...v4.1.1
[v4.1.0]: https://github.com/beyondstorage/go-integration-test/compare/v4.0.0...v4.1.0
[v4.0.0]: https://github.com/beyondstorage/go-integration-test/compare/v3.0.0...v4.0.0
