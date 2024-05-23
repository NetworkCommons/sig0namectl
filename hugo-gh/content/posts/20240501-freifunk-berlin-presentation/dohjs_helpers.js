/// DoHjs Helper Functions
///
/// These functions help getting the values from binary wireformat
/// which are not yet handled by DoHjs.

/// Uint32 helper Class for byte conversion
class Uint32
{
    constructor(Value) 
    {
        this.Number = new Uint32Array(1);
        this.Number[0] = Value;
    }
    get Get() 
    {
        return this.Number[0];
    }
    set Set(newValue) 
    {
        this.Number[0] = newValue;
    }
};

/// Convert u8 Byte array from network byte order (big-endian) to 32bit unsigned integer (little-endian)
function ntoh_Uint32 (Source_Byte_Array, Start_Position)
{
    var Uint32_Num = new Uint32(0);
    var Multiplier = 1;
    for (let i = 3; i >= 0; i--)
    {
        Uint32_Num.Set = Uint32_Num.Get + Source_Byte_Array[Start_Position + i] * Multiplier;
        Multiplier = Multiplier * 256;
    }
    return (Uint32_Num.Get);
}

/// LOC Wireformat Decoder
class LocDecoder
{
    /// constructor takes the raw wireformat bytes array
    constructor(wire_data) 
    {
        this.data = wire_data;
        this.latitude = this.get_latitude();
        this.longitude = this.get_longitude();
    }

    /// get latitude
    get_latitude() {
        let latitude_raw = ntoh_Uint32(this.data, 4);
        let latitude = (latitude_raw - 2147483648.0) / 3600000;
        return latitude;
    }

    /// get longitude
    get_longitude() {
        let longitude_raw = ntoh_Uint32(this.data, 8);
        let longitude = (longitude_raw - 2147483648.0) / 3600000;
        return longitude;
    }

}