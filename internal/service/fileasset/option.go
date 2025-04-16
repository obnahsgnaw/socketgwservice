package fileasset

type Option func(*Server)

func MaxSize(size uint32) Option {
	return func(s *Server) {
		s.maxSize = size
	}
}

func MaxCount(count uint32) Option {
	return func(s *Server) {
		s.maxCount = count
	}
}

func Multipart() Option {
	return func(s *Server) {
		s.multipart = true
	}
}

func ViewUrlTtl(ttl uint32) Option {
	return func(s *Server) {
		s.viewUrlTtl = ttl
	}
}

func UploadUrlTtl(ttl uint32) Option {
	return func(s *Server) {
		s.UploadUrlTtl = ttl
	}
}

func SessionIdTtl(ttl uint32) Option {
	return func(s *Server) {
		s.sessionTtl = ttl
	}
}

func Enable(cb func() bool) Option {
	return func(s *Server) {
		s.enableChecker = cb
	}
}
