class UrlModel {
  final String originalUrl;
  final String shortUrl;
  final String code;
  final DateTime createdAt;
  final DateTime? lastAccessed;
  final int clickCount;

  UrlModel({
    required this.originalUrl,
    required this.shortUrl,
    required this.code,
    required this.createdAt,
    this.lastAccessed,
    required this.clickCount,
  });

  factory UrlModel.fromJson(Map<String, dynamic> json) {
    return UrlModel(
      originalUrl: json['original_url'] ?? '',
      shortUrl: json['short_url'] ?? '',
      code: json['code'] ?? '',
      createdAt: DateTime.parse(json['created_at']),
      lastAccessed: json['last_accessed'] != null 
          ? DateTime.parse(json['last_accessed'])
          : null,
      clickCount: json['click_count'] ?? 0,
    );
  }
}

class ShortenRequest {
  final String url;
  final String? customCode;

  ShortenRequest({required this.url, this.customCode});

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = {'url': url};
    if (customCode != null && customCode!.isNotEmpty) {
      data['custom_code'] = customCode;
    }
    return data;
  }
}
