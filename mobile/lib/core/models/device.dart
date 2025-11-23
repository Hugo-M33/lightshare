/// Device model representing a smart light device from any provider
class Device {
  final String id;
  final String accountId;
  final String provider;
  final String label;
  final String power; // "on" or "off"
  final double brightness; // 0.0 - 1.0
  final DeviceColor? color;
  final bool connected;
  final bool reachable;
  final DeviceGroup? group;
  final DeviceLocation? location;
  final List<String> capabilities;
  final Map<String, dynamic>? metadata;

  Device({
    required this.id,
    required this.accountId,
    required this.provider,
    required this.label,
    required this.power,
    required this.brightness,
    this.color,
    required this.connected,
    required this.reachable,
    this.group,
    this.location,
    required this.capabilities,
    this.metadata,
  });

  factory Device.fromJson(Map<String, dynamic> json) {
    return Device(
      id: json['id'] as String,
      accountId: json['account_id'] as String,
      provider: json['provider'] as String,
      label: json['label'] as String,
      power: json['power'] as String,
      brightness: (json['brightness'] as num).toDouble(),
      color: json['color'] != null
          ? DeviceColor.fromJson(json['color'] as Map<String, dynamic>)
          : null,
      connected: json['connected'] as bool,
      reachable: json['reachable'] as bool,
      group: json['group'] != null
          ? DeviceGroup.fromJson(json['group'] as Map<String, dynamic>)
          : null,
      location: json['location'] != null
          ? DeviceLocation.fromJson(json['location'] as Map<String, dynamic>)
          : null,
      capabilities: (json['capabilities'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      metadata: json['metadata'] as Map<String, dynamic>?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'account_id': accountId,
      'provider': provider,
      'label': label,
      'power': power,
      'brightness': brightness,
      'color': color?.toJson(),
      'connected': connected,
      'reachable': reachable,
      'group': group?.toJson(),
      'location': location?.toJson(),
      'capabilities': capabilities,
      'metadata': metadata,
    };
  }

  Device copyWith({
    String? id,
    String? accountId,
    String? provider,
    String? label,
    String? power,
    double? brightness,
    DeviceColor? color,
    bool? connected,
    bool? reachable,
    DeviceGroup? group,
    DeviceLocation? location,
    List<String>? capabilities,
    Map<String, dynamic>? metadata,
  }) {
    return Device(
      id: id ?? this.id,
      accountId: accountId ?? this.accountId,
      provider: provider ?? this.provider,
      label: label ?? this.label,
      power: power ?? this.power,
      brightness: brightness ?? this.brightness,
      color: color ?? this.color,
      connected: connected ?? this.connected,
      reachable: reachable ?? this.reachable,
      group: group ?? this.group,
      location: location ?? this.location,
      capabilities: capabilities ?? this.capabilities,
      metadata: metadata ?? this.metadata,
    );
  }

  // Helper getters
  bool get isOn => power == 'on';

  bool get isOff => power == 'off';

  bool get isAvailable => connected && reachable;

  bool hasCapability(String capability) => capabilities.contains(capability);

  bool get supportsColor => hasCapability('color');

  bool get supportsTemperature => hasCapability('temperature');

  bool get supportsEffects => hasCapability('effects');

  String get displayName => label.isNotEmpty ? label : 'Device $id';

  String get providerDisplayName => provider.toUpperCase();
}

/// Device color information
class DeviceColor {
  final double hue; // 0-360 degrees
  final double saturation; // 0.0-1.0
  final int kelvin; // 1500-9000

  DeviceColor({
    required this.hue,
    required this.saturation,
    required this.kelvin,
  });

  factory DeviceColor.fromJson(Map<String, dynamic> json) {
    return DeviceColor(
      hue: (json['hue'] as num).toDouble(),
      saturation: (json['saturation'] as num).toDouble(),
      kelvin: json['kelvin'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'hue': hue,
      'saturation': saturation,
      'kelvin': kelvin,
    };
  }

  DeviceColor copyWith({
    double? hue,
    double? saturation,
    int? kelvin,
  }) {
    return DeviceColor(
      hue: hue ?? this.hue,
      saturation: saturation ?? this.saturation,
      kelvin: kelvin ?? this.kelvin,
    );
  }
}

/// Device group/room information
class DeviceGroup {
  final String id;
  final String name;

  DeviceGroup({
    required this.id,
    required this.name,
  });

  factory DeviceGroup.fromJson(Map<String, dynamic> json) {
    return DeviceGroup(
      id: json['id'] as String,
      name: json['name'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
    };
  }

  DeviceGroup copyWith({
    String? id,
    String? name,
  }) {
    return DeviceGroup(
      id: id ?? this.id,
      name: name ?? this.name,
    );
  }
}

/// Device location/home information
class DeviceLocation {
  final String id;
  final String name;

  DeviceLocation({
    required this.id,
    required this.name,
  });

  factory DeviceLocation.fromJson(Map<String, dynamic> json) {
    return DeviceLocation(
      id: json['id'] as String,
      name: json['name'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
    };
  }

  DeviceLocation copyWith({
    String? id,
    String? name,
  }) {
    return DeviceLocation(
      id: id ?? this.id,
      name: name ?? this.name,
    );
  }
}
